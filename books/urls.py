from django.conf import settings
from django.conf.urls.static import static
from django.urls import path
from .views import *

urlpatterns = [
    path("category/create", CategoriesCreateAPIView.as_view()),
    path("category/all", CategoriesListAPIView.as_view()),
    path("category/<int:pk>", CategoriesRetrieveUpdateDeleteAPIView.as_view()),
    path("category/create", CategoriesCreateAPIView.as_view()),
    path("category/all", CategoriesListAPIView.as_view()),
    path("category/<int:pk>", CategoriesRetrieveUpdateDeleteAPIView.as_view()),
    path("book/create", BooksCreateAPIView.as_view()),
    path("book/all", BooksListAPIView.as_view()),
    path("book/<int:pk>", BooksRetrieveUpdateDeleteAPIView.as_view()),
    path("book/<int:pk>", EBookDownloadView.as_view())
] + static(settings.MEDIA_URL, document_root=settings.MEDIA_ROOT)
