from rest_framework.generics import CreateAPIView, ListAPIView, RetrieveUpdateDestroyAPIView
from rest_framework.permissions import IsAuthenticated, AllowAny
from rest_framework.response import Response
from rest_framework.status import HTTP_403_FORBIDDEN, HTTP_404_NOT_FOUND, HTTP_201_CREATED, HTTP_200_OK
from django.db.models import Q
from django.core.files.storage import default_storage
from users.permissions import IsLibrarian
from core.settings import ERROR_404_IMAGE_FOLDER
from .models import Book, Category
from .serializers import BookSerializer, CategorySerializer


class CategoriesCreateAPIView(CreateAPIView):
    queryset = Category.objects.all()
    serializer_class = CategorySerializer
    permission_classes = [IsAuthenticated, IsLibrarian]


class CategoriesListAPIView(ListAPIView):
    queryset = Category.objects.all()
    serializer_class = CategorySerializer
    permission_classes = [AllowAny]


class CategoriesRetrieveUpdateDeleteAPIView(RetrieveUpdateDestroyAPIView):
    queryset = Category.objects.all()
    serializer_class = CategorySerializer
    permission_classes = [IsAuthenticated, IsLibrarian]


class BooksCreateAPIView(CreateAPIView):
    queryset = Book.objects.all()
    serializer_class = BookSerializer
    permission_classes = [IsAuthenticated, IsLibrarian]

    def post(self, request, *args, **kwargs):
        serializer = BookSerializer(data=request.data)
        serializer.is_valid(raise_exception=True)

        if serializer.validated_data["quantity"] == 0:
            serializer.validated_data["isPossibleToOrder"] = False
        if serializer.validated_data["quantity"] < 0:
            raise ValueError("Количество книг не может быть отрицательным")
        serializer.save()

        return Response(
            {"message": "Книга успешно добавлена"}, status=HTTP_201_CREATED
        )


class BooksListAPIView(ListAPIView):
    queryset = Book.objects.all()
    serializer_class = BookSerializer
    permission_classes = [AllowAny]

    def get(self, request, *args, **kwargs):
        category = request.GET.get("category")
        less_orders = request.GET.get("less_orders")
        more_orders = request.GET.get("more_orders")
        author = request.GET.get("author")
        title = request.GET.get("title")

        for book in self.get_queryset().all():
            image_path = book.image.name
            if default_storage.exists(image_path):
                default_storage.url(image_path)
                continue
            book.image = ERROR_404_IMAGE_FOLDER
            book.save()

        if category:
            category = category.capitalize()
            self.queryset = self.queryset.filter(category__title=category)
        if less_orders:
            self.queryset = self.queryset.filter(orders__lte=less_orders)
        if more_orders:
            self.queryset = self.queryset.filter(orders__gte=more_orders)

        if author:
            self.queryset = self.queryset.filter(Q(author__icontains=author))
        elif title:
            self.queryset = self.queryset.filter(Q(title__icontains=title))

        return super().get(request)


class BooksRetrieveUpdateDeleteAPIView(RetrieveUpdateDestroyAPIView):
    queryset = Book.objects.all()
    serializer_class = BookSerializer
    permission_classes = [IsAuthenticated]

    def get(self, request, *args, **kwargs):
        book = self.get_object()

        if not book:
            return Response({"message": "Книга не найдена"}, status=HTTP_404_NOT_FOUND)

        image_path = book.image.name
        if default_storage.exists(image_path):
            default_storage.url(image_path)
            return super().get(request, *args, **kwargs)
        book.image = ERROR_404_IMAGE_FOLDER
        book.save()
        return super().get(request, *args, **kwargs)

    def put(self, request, *args, **kwargs):
        if self.request.user.status in ["Student", "Admin"]:
            return Response({"message": "Вы не можете изменить книгу"}, status=HTTP_403_FORBIDDEN)

        book = self.get_object()

        if not book:
            return Response({"message": "Книга не найдена"}, status=HTTP_404_NOT_FOUND)

        if book.image != request.data["image"] \
                and book.image.path != ERROR_404_IMAGE_FOLDER:
            default_storage.delete(book.image.path)

        return Response({"message": "Книга успешно изменена"})

    def delete(self, request, *args, **kwargs):
        if self.request.user.status in ["Student", "Admin"]:
            return Response({"message": "Вы не можете удалить книгу"}, status=HTTP_403_FORBIDDEN)

        book = self.get_object()

        if not book:
            return Response({"message": "Книга не найдена"}, status=HTTP_404_NOT_FOUND)

        if book.image and book.image.path != ERROR_404_IMAGE_FOLDER:
            default_storage.delete(book.image.path)
        book.delete()
        return Response({"message": "Книга успешно удалена"}, status=HTTP_200_OK)
