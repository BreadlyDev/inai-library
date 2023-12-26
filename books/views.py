from datetime import datetime
from django.http import FileResponse
from rest_framework.generics import CreateAPIView, ListAPIView, RetrieveUpdateDestroyAPIView, RetrieveAPIView
from rest_framework.permissions import IsAuthenticated, AllowAny
from rest_framework.response import Response
from rest_framework.status import HTTP_403_FORBIDDEN, HTTP_404_NOT_FOUND, HTTP_201_CREATED, HTTP_200_OK
from rest_framework.views import APIView
from django.db.models import Q
from django.core.files.storage import default_storage
from users.permissions import IsLibrarian
from main.settings import ERROR_404_IMAGE, REPORTS_FOLDER, BASE_DIR
from .models import Book, Category, Subcategory
from .serializers import BookSerializer, CategorySerializer, SubcategorySerializer
from .reports import create_report, fill_table


class SubcategoriesCreateAPIView(CreateAPIView):
    queryset = Subcategory.objects.all()
    serializer_class = SubcategorySerializer
    permission_classes = [IsAuthenticated, IsLibrarian]


class SubcategoriesListAPIView(ListAPIView):
    queryset = Subcategory.objects.all()
    serializer_class = SubcategorySerializer
    permission_classes = [AllowAny]


class SubcategoriesRetrieveUpdateDeleteAPIView(RetrieveUpdateDestroyAPIView):
    queryset = Subcategory.objects.all()
    serializer_class = SubcategorySerializer
    permission_classes = [IsAuthenticated, IsLibrarian]


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
            if not book.image:
                book.image = ERROR_404_IMAGE
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
        book.image = ERROR_404_IMAGE
        book.save()
        return super().get(request, *args, **kwargs)

    def put(self, request, *args, **kwargs):
        if self.request.user.status in ["Student", "Admin"]:
            return Response({"message": "Вы не можете изменить книгу"}, status=HTTP_403_FORBIDDEN)

        book = self.get_object()

        if not book:
            return Response({"message": "Книга не найдена"}, status=HTTP_404_NOT_FOUND)

        if book.image != request.data["image"] \
                and book.image.path != ERROR_404_IMAGE:
            default_storage.delete(book.image.path)

        return Response({"message": "Книга успешно изменена"})

    def delete(self, request, *args, **kwargs):
        if self.request.user.status in ["Student", "Admin"]:
            return Response({"message": "Вы не можете удалить книгу"}, status=HTTP_403_FORBIDDEN)

        book = self.get_object()

        if not book:
            return Response({"message": "Книга не найдена"}, status=HTTP_404_NOT_FOUND)

        if book.image and book.image.path != ERROR_404_IMAGE:
            default_storage.delete(book.image.path)
        book.delete()
        return Response({"message": "Книга успешно удалена"}, status=HTTP_200_OK)


class EBookDownloadView(RetrieveAPIView):
    queryset = Book.objects.all()
    serializer_class = BookSerializer

    def retrieve(self, request, *args, **kwargs):
        instance = self.get_object()
        file_path = instance.e_book.path

        if not file_path:
            return Response({"message": "Файл отсутствует"})

        response = FileResponse(open(file_path, "rb"))
        response["Content-Disposition"] = f"attachment; filename={instance.file_field.name}"
        return response


class BookReportCreateAPIView(CreateAPIView):
    permission_classes = [IsAuthenticated, IsLibrarian]

    def post(self, request, *args, **kwargs):
        try:
            subcategories = Subcategory.objects.all()
            books = Book.objects.all()
            document, table = create_report()

            for subcategory in subcategories:
                j = 0
                for book in books:
                    if book.subcategory == subcategory:
                        row_cells = table.add_row().cells
                        table_texts = ["",
                                       str(subcategory.title) if j == 0 else "",
                                       "Очная/компьютерные  технологии" if j == 0 else "",
                                       str(book.inventory_number),
                                       str(book.quantity),
                                       str(book.author),
                                       str(book.title),
                                       str(book.edition_year)]
                        fill_table(row_cells=row_cells, table_texts=table_texts)
                        j += 1

            document.save(
                f"{BASE_DIR}/media/{REPORTS_FOLDER}отчёт_за_{datetime.now().strftime('%d-%m-%Y_%H-%M-%S')}.docx")
            print(f"{BASE_DIR}/media/{REPORTS_FOLDER}отчёт_за_{datetime.now().strftime('%d-%m-%Y_%H-%M-%S')}.docx")
            return Response({"message": "Report created successfully"})
        except PermissionError:
            return Response({"message": "You should first close the file before creating it again"})


# class BookReportCreateAPIView(ListAPIView):
#     queryset = Book.objects.all()
#     serializer_class = BookSerializer
#
#     def retrieve(self, request, *args, **kwargs):
#         instance = self.get_object()
#         file_path = instance.e_book.path
#
#         if not file_path:
#             return Response({"message": "Файл отсутствует"})
#
#         response = FileResponse(open(file_path, "rb"))
#         response["Content-Disposition"] = f"attachment; filename={instance.file_field.name}"
#         return response
