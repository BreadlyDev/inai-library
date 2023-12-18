from rest_framework.permissions import BasePermission


class IsLibrarian(BasePermission):
    def has_permission(self, request, view):
        return request.user.role == "Librarian"


class IsAdmin(BasePermission):
    def has_permission(self, request, view):
        return request.user.role == "Admin"


class IsStudent(BasePermission):
    def has_permission(self, request, view):
        return request.user.role == "Student"


class IsLibrarianOrStudent(BasePermission):
    def has_permission(self, request, view):
        return request.user.role == "Student" or "Librarian"


class NotStudentPermission(BasePermission):
    def has_permission(self, request, view):
        return not request.user.role == "Student"
