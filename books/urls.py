from django.conf import settings
from django.conf.urls.static import static
from django.urls import path
from .views import *

urlpatterns = [
    path("create/category", CategoriesCreateAPIView.as_view()),
    path("list/category", CategoriesListAPIView.as_view()),
    path("change/category/<int:pk>", CategoriesRetrieveUpdateDeleteAPIView.as_view()),
    path("create/book", BooksCreateAPIView.as_view()),
    path("list/book", BooksListAPIView.as_view()),
    path("change/book/<int:pk>", BooksRetrieveUpdateDeleteAPIView.as_view()),
] + static(settings.MEDIA_URL, document_root=settings.MEDIA_ROOT)
