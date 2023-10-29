from rest_framework.response import Response
from rest_framework_simplejwt.tokens import RefreshToken
from rest_framework.permissions import AllowAny
from rest_framework.generics import CreateAPIView
from rest_framework.status import HTTP_201_CREATED, HTTP_200_OK
from rest_framework.views import APIView
from django.contrib.auth import logout
from .serializers import UserSerializer, LoginSerializer
from .models import User


class UserRegisterAPIView(CreateAPIView):
    queryset = User.objects.all()
    serializer_class = UserSerializer
    permission_classes = (AllowAny,)

    def post(self, request, *args, **kwargs):
        data = request.data
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
                'refresh_token': refresh_token
            },
            status=HTTP_201_CREATED
        )


class UserLoginAPIView(APIView):
    queryset = User.objects.all()
    serializer_class = LoginSerializer
    permission_classes = (AllowAny,)

    @staticmethod
    def post(request):
        data = request.data
        email = data.get("email")
        password = data.get("password")
        user = User.objects.filter(email=email).first()

        if user is None:
            return Response({'error': 'User not found'}, status=400)

        if not user.check_password(password):
            return Response({'error': 'Invalid password'}, status=400)

        refresh = RefreshToken.for_user(user)
        access_token = str(refresh.access_token)
        refresh_token = str(refresh)

        return Response(
            {
                'message': 'User logged in successfully',
                'access_token': access_token,
                'refresh_token': refresh_token
            },
            status=HTTP_200_OK
        )


class UserLogoutAPIView(APIView):
    @classmethod
    def post(request):
        logout(request)
        return Response({'message': 'User logged out successfully'}, status=HTTP_200_OK)
