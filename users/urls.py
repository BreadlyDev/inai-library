from django.urls import path
from rest_framework_simplejwt.views import TokenRefreshView
from .views import (UserRegisterAPIView, UserLoginAPIView, UserLogoutAPIView, UserListAPIView, UserGetAPIView,
                    GroupCreateAPIView, GroupChangeAPIView, GroupListAPIView)

urlpatterns = [
    path('register/', UserRegisterAPIView.as_view()),
    path('login/', UserLoginAPIView.as_view()),
    path('logout/', UserLogoutAPIView.as_view()),
    path('get/user/', UserGetAPIView.as_view()),
    path('list/user/', UserListAPIView.as_view()),
    path('create/group/', GroupCreateAPIView.as_view()),
    path('change/group/<int:pk>/', GroupChangeAPIView.as_view()),
    path('list/group/', GroupListAPIView.as_view()),
    path('activate/refresh/token/', TokenRefreshView.as_view()),
]
