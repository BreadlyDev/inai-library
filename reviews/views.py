from rest_framework.response import Response
from rest_framework.status import HTTP_201_CREATED, HTTP_403_FORBIDDEN, HTTP_204_NO_CONTENT

from .models import Review
from .serializers import ReviewSerializer
from rest_framework.generics import ListAPIView, CreateAPIView, RetrieveUpdateDestroyAPIView
from accounts.permissions import IsStudent
from rest_framework.permissions import IsAuthenticated, AllowAny


class ReviewListAPIView(ListAPIView):
    queryset = Review.objects.all()
    serializer_class = ReviewSerializer
    permission_classes = [AllowAny]


class ReviewCreateAPIView(CreateAPIView):
    queryset = Review.objects.all()
    serializer_class = ReviewSerializer
    permission_classes = [IsAuthenticated, IsStudent]

    def create(self, request, *args, **kwargs):
        serializer = ReviewSerializer(data=request.data)
        serializer.is_valid(raise_exception=True)
        serializer.validated_data["author"] = request.user
        book = serializer.validated_data["book"]
        serializer.save()

        book.reviews += 1
        book.save()

        return Response(serializer.data, status=HTTP_201_CREATED)


class ReviewChangeAPIView(RetrieveUpdateDestroyAPIView):
    queryset = Review.objects.all()
    serializer_class = ReviewSerializer
    permission_classes = [IsAuthenticated, IsStudent]

    def update(self, request, *args, **kwargs):
        review = self.get_object()
        if review.author == request.user:
            return super().update(request, *args, **kwargs)
        else:
            return Response(
                {'message': 'У вас нет разрешения на изменение этого отзыва'},
                status=HTTP_403_FORBIDDEN,
            )

    def delete(self, request, *args, **kwargs):
        review = self.get_object()
        if review.author == request.user:
            review.delete()
            return Response(
                {'message': 'Отзыв успешно удален'}, status=HTTP_204_NO_CONTENT
            )
        else:
            return Response(
                {'message': 'У вас нет разрешения на удаление этого отзыва'},
                status=HTTP_403_FORBIDDEN
            )
