from rest_framework.response import Response
from rest_framework_simplejwt.tokens import RefreshToken
from rest_framework.permissions import AllowAny, IsAuthenticated, IsAuthenticatedOrReadOnly
from rest_framework.generics import CreateAPIView, ListAPIView, RetrieveUpdateDestroyAPIView
from rest_framework.status import HTTP_201_CREATED, HTTP_200_OK, HTTP_400_BAD_REQUEST
from rest_framework.views import APIView
from .permissions import NotStudentPermission, IsAdmin, IsLibrarian
from .serializers import UserSerializer, LoginSerializer, GroupSerializer
from .models import User, Group


class UserRegisterAPIView(CreateAPIView):
    queryset = User.objects.all()
    serializer_class = UserSerializer
    permission_classes = [AllowAny]

    def post(self, request, *args, **kwargs):
        data = request.data

        if not data["phone"][1:].isdigit() \
                or not 8 < len(data["phone"]) < 14:
            return Response({"message": "Неверный номер телефона"}, status=HTTP_400_BAD_REQUEST)

        serializer = self.get_serializer(data=data)
        serializer.is_valid(raise_exception=True)
        user = serializer.save()

        refresh = RefreshToken.for_user(user)
        access_token = str(refresh.access_token)
        refresh_token = str(refresh)

        return Response(
            {
                'message': 'User registered successfully',
                'access_token': access_token,
                'refresh_token': refresh_token,
                'user': serializer.data,
            },
            status=HTTP_201_CREATED
        )


class UserLoginAPIView(APIView):
    queryset = User.objects.all()
    serializer_class = LoginSerializer
    permission_classes = [AllowAny]

    @staticmethod
    def post(request):
        data = request.data
        email = data.get("email")
        password = data.get("password")
        user = User.objects.filter(email=email).first()
        serializer = UserSerializer(user)

        if user is None:
            return Response({'message': 'User not found'}, status=400)

        if not user.check_password(password):
            return Response({'message': 'Invalid password'}, status=400)

        refresh = RefreshToken.for_user(user)
        access_token = str(refresh.access_token)
        refresh_token = str(refresh)

        return Response(
            {
                'message': 'User logged in successfully',
                'access_token': access_token,
                'refresh_token': refresh_token,
                'user': serializer.data,
            },
            status=HTTP_200_OK
        )


class UserLogoutAPIView(APIView):

    @classmethod
    def post(cls, request):
        refresh_token = request.data.get('refresh_token')

        if not refresh_token:
            return Response({'message': "Отсутствует Refresh токен"}, status=HTTP_400_BAD_REQUEST)

        try:
            RefreshToken(refresh_token).blacklist()
            return Response({'message': 'Пользователь успещно вышел из системы.'}, status=HTTP_200_OK)
        except Exception:
            return Response({'message': 'Неверный токен или токен просрочен.'}, status=HTTP_400_BAD_REQUEST)


class UserListAPIView(ListAPIView):
    queryset = User.objects.filter(is_superuser=False)
    serializer_class = UserSerializer
    permission_classes = [IsAuthenticated, NotStudentPermission]


class UserGetAPIView(APIView):
    permission_classes = [IsAuthenticated, IsLibrarian]

    @classmethod
    def get(cls, request, pk):
        user = User.objects.filter(pk=pk)
        # user = request.user
        serialized_user = UserSerializer(user).to_representation(user)
        return Response(serialized_user.data)


class GroupCreateAPIView(CreateAPIView):
    serializer_class = GroupSerializer
    permission_classes = [IsAuthenticated, IsAdmin]


class GroupChangeAPIView(RetrieveUpdateDestroyAPIView):
    queryset = Group.objects.all()
    serializer_class = GroupSerializer
    permission_classes = [IsAuthenticatedOrReadOnly, IsAdmin]


class GroupListAPIView(ListAPIView):
    queryset = Group.objects.all()
    serializer_class = GroupSerializer
    permission_classes = [IsAuthenticated]
