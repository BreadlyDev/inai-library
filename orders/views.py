from rest_framework.response import Response
from rest_framework.generics import ListAPIView, CreateAPIView, RetrieveUpdateDestroyAPIView
from rest_framework.permissions import IsAuthenticated
from rest_framework.status import HTTP_403_FORBIDDEN, HTTP_204_NO_CONTENT, HTTP_201_CREATED
from users.permissions import IsStudent, IsLibrarianOrStudent
from .models import Order
from .serializers import OrderCreateSerializer, OrderSerializer


class OrderCreateAPIView(CreateAPIView):
    queryset = Order.objects.all()
    serializer_class = OrderCreateSerializer
    permission_classes = [IsAuthenticated, IsStudent]

    def post(self, request, *args, **kwargs):
        serializer = OrderCreateSerializer(data=request.data)
        serializer.is_valid(raise_exception=True)
        serializer.validated_data["owner"] = request.user
        books = serializer.validated_data["books"]

        for book in books:
            if book.quantity <= 0:
                book.isPossibleToOrder = False
            if not book.isPossibleToOrder:
                return Response({"message": "К сожалению вы не можете забронировать эту книгу на данный момент"})

        order = serializer.save()

        for book in order.books.all():
            book.orders += 1
            book.quantity -= 1
            book.save()

        return Response(serializer.data, status=HTTP_201_CREATED)


class OrderListAPIView(ListAPIView):
    queryset = Order.objects.all()
    serializer_class = OrderSerializer
    permission_classes = [IsAuthenticated]

    def get_queryset(self):
        user = self.request.user

        if user.status == "Student":
            return Order.objects.filter(owner=user)

        return Order.objects.all()


class OrderRetrieveUpdateDestroyAPIView(RetrieveUpdateDestroyAPIView):
    queryset = Order.objects.all()
    serializer_class = OrderCreateSerializer
    permission_classes = [IsAuthenticated, IsLibrarianOrStudent]

    def get_serializer_class(self):
        if self.request.method == "GET":
            return OrderSerializer

    def put(self, request, *args, **kwargs):
        order = self.get_object()

        if (
                request.user.status == "Student"
                and order.status == "Не рассмотрено"
                and order.owner == request.user
        ):
            return super().put(request, *args, **kwargs)

        return Response(
                {'message': 'У вас нет разрешения на изменение этого заказа'},
                status=HTTP_403_FORBIDDEN,
        )

    def delete(self, request, *args, **kwargs):
        order = self.get_object()

        if (
            request.user.status == "Student"
            and order.status == "Не рассмотрено"
            and order.owner == request.user
        ):
            order.delete()
            return Response(
                {'message': 'Заказ удален успешно'}, status=HTTP_204_NO_CONTENT
            )

        return Response(
                {'message': 'У вас нет разрешения на удаление этого заказа'}, status=HTTP_403_FORBIDDEN
        )
