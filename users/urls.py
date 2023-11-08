from django.urls import path
from .views import UserRegisterAPIView, UserLoginAPIView, UserLogoutAPIView, UserListAPIView, UserGetAPIView

urlpatterns = [
    path('register/', UserRegisterAPIView.as_view()),
    path('login/', UserLoginAPIView.as_view()),
    path('logout/', UserLogoutAPIView.as_view()),
    path('get/user/', UserGetAPIView.as_view()),
    path('list/user/', UserListAPIView.as_view()),
]
