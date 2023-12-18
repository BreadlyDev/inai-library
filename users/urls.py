from django.urls import path
from rest_framework_simplejwt.views import TokenRefreshView
from .views import (UserRegisterAPIView, UserLoginAPIView, UserLogoutAPIView, UserListAPIView, UserGetAPIView,
                    GroupCreateAPIView, GroupChangeAPIView, GroupListAPIView)

urlpatterns = [
    path('register', UserRegisterAPIView.as_view()),
    path('login', UserLoginAPIView.as_view()),
    path('logout', UserLogoutAPIView.as_view()),
    path('user/<int:pk>', UserGetAPIView.as_view()),
    path('user/all', UserListAPIView.as_view()),
    path('group/create', GroupCreateAPIView.as_view()),
    path('group/<int:pk>', GroupChangeAPIView.as_view()),
    path('group/all', GroupListAPIView.as_view()),
    path('activate/refresh/token', TokenRefreshView.as_view()),
]
