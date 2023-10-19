from rest_framework.response import Response
from .models import Order
from .serializers import OrderSerializer
from accounts.permissions import IsStudent, IsLibrarianOrStudent
from rest_framework.generics import ListAPIView, CreateAPIView, RetrieveUpdateDestroyAPIView
from rest_framework.permissions import IsAuthenticated

from rest_framework.status import HTTP_403_FORBIDDEN, HTTP_204_NO_CONTENT, HTTP_201_CREATED


class OrderCreateAPIView(CreateAPIView):
    queryset = Order.objects.all()
    serializer_class = OrderSerializer
    permission_classes = [IsAuthenticated, IsStudent]

    def create(self, request, *args, **kwargs):
        serializer = OrderSerializer(data=request.data)
        serializer.is_valid(raise_exception=True)
        serializer.validated_data["owner"] = request.user

        if not serializer.validated_data["isPossibleToOrder"]:
            return Response({"message" : "К сожалению вы не можете заюроинровать эту книгу на данный момент"})


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


class OrderChangeAPIView(RetrieveUpdateDestroyAPIView):
    queryset = Order.objects.all()
    serializer_class = OrderSerializer
    permission_classes = [IsAuthenticated, IsLibrarianOrStudent]

    def update(self, request, *args, **kwargs):
        order = self.get_object()

        if(
                request.user.status == "Student"
                and order.status == "Не рассмотрено"
                and order.owner == request.user
        ):
            return super().update(request, *args, **kwargs)
        else:
            return Response(
                {'message': 'У вас нет разрешения на изменение этого заказа'},
                status=HTTP_403_FORBIDDEN,
            )

    def delete(self, request, *args, **kwargs):
        order = self.get_object()

        if(
                request.user.status == "Student"
                and order.status == "Не рассмотрено"
                and order.owner == request.user
        ):
            order.delete()
            return Response(
                {'message': 'Заказ удален успешно'}, status=HTTP_204_NO_CONTENT
            )
        else:
            return Response(
                {'message': 'У вас нет разрешения на удаление этого заказа'}, status=HTTP_403_FORBIDDEN
            )
