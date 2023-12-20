from rest_framework.permissions import BasePermission, SAFE_METHODS
from users.models import ROLES


class IsAdminOrReadOnly(BasePermission):
    def has_permission(self, request, view):
        if request.method in SAFE_METHODS:
            return True
        return request.user.is_authenticated and request.user.role == ROLES[0][1]


class IsLibrarian(BasePermission):
    def has_permission(self, request, view):
        if request.user.is_authenticated:
            return request.user.role == ROLES[1][1]
        return False


class IsAdmin(BasePermission):
    def has_permission(self, request, view):
        if request.user.is_authenticated:
            return request.user.role == ROLES[0][1]
        return False


class IsStudent(BasePermission):
    def has_permission(self, request, view):
        if request.user.is_authenticated:
            return request.user.role == ROLES[2][1]
        return False


class IsLibrarianOrStudent(BasePermission):
    def has_permission(self, request, view):
        if request.user.is_authenticated:
            return request.user.role == ROLES[1][1] or ROLES[2][1]
        return False


class NotStudentPermission(BasePermission):
    def has_permission(self, request, view):
        return not request.user.role == ROLES[2][1]
