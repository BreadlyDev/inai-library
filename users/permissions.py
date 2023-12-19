from rest_framework.permissions import BasePermission, SAFE_METHODS


class IsAdminOrReadOnly(BasePermission):
    def has_permission(self, request, view):
        if request.method in SAFE_METHODS:
            return True
        return request.user.is_authenticated and request.user.role == "Admin"


class IsLibrarian(BasePermission):
    def has_permission(self, request, view):
        if request.user.is_authenticated:
            return request.user.role == "Librarian"
        return False


class IsAdmin(BasePermission):
    def has_permission(self, request, view):
        if request.user.is_authenticated:
            return request.user.role == "Admin"
        return False


class IsStudent(BasePermission):
    def has_permission(self, request, view):
        if request.user.is_authenticated:
            return request.user.role == "Student"
        return False


class IsLibrarianOrStudent(BasePermission):
    def has_permission(self, request, view):
        if request.user.is_authenticated:
            return request.user.role == "Student" or "Librarian"
        return False


class NotStudentPermission(BasePermission):
    def has_permission(self, request, view):
        return not request.user.role == "Student"
