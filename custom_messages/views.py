from rest_framework.generics import ListAPIView, CreateAPIView, RetrieveUpdateDestroyAPIView
from rest_framework.response import Response
from rest_framework.status import HTTP_201_CREATED, HTTP_403_FORBIDDEN
from rest_framework.permissions import IsAuthenticated
from users.permissions import IsLibrarian, IsLibrarianOrStudent
from .models import Message
from .serializers import MessageSerializer


class MessageListAPIView(ListAPIView):
    queryset = Message.objects.all()
    serializer_class = MessageSerializer
    permission_classes = [IsAuthenticated, IsLibrarianOrStudent]

    def get_queryset(self):
        user = self.request.user

        if user.status == "Student":
            return Message.objects.filter(recipient=user)
        return Message.objects.all()


class MessageCreateAPIView(CreateAPIView):
    queryset = Message.objects.all()
    serializer_class = MessageSerializer
    permission_classes = [IsAuthenticated, IsLibrarian]

    def post(self, request, *args, **kwargs):
        serializer = MessageSerializer(data=request.data)
        serializer.is_valid(raise_exception=True)
        recipient = serializer.validated_data.get("recipient")
        if recipient.status == "Student":
            serializer.save()
            return Response(serializer.data, status=HTTP_201_CREATED)

        return Response(
            {'message': 'Вы можете отправлять сообщения только студентам.'},
            status=HTTP_403_FORBIDDEN
        )


class MessageRetrieveUpdateDeleteAPIView(RetrieveUpdateDestroyAPIView):
    queryset = Message.objects.all()
    serializer_class = MessageSerializer
    permission_classes = [IsAuthenticated, IsLibrarian]
