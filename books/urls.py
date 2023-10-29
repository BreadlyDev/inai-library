from django.urls import path
from .views import *

urlpatterns = [
    path("create/category/", CategoriesCreateAPIView.as_view()),
    path("list/category/", CategoriesListAPIView.as_view()),
    path("retrieve/update/delete/category/<int:pk>/", CategoriesRetrieveUpdateDeleteAPIView.as_view()),
    path("create/book/", BooksCreateAPIView.as_view()),
    path("list/book/", BooksListAPIView.as_view()),
    path("retrieve/update/delete/book/<int:pk>/", BooksRetrieveUpdateDeleteAPIView.as_view()),
]
