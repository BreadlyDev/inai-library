from rest_framework.permissions import BasePermission


class IsLibrarian(BasePermission):
    def has_permission(self, request, view):
        return request.user.status == "Librarian"


class IsAdmin(BasePermission):
    def has_permission(self, request, view):
        return request.user.status == "Admin"


class IsStudent(BasePermission):
    def has_permission(self, request, view):
        return request.user.status == "Student"


class IsLibrarianOrStudent(BasePermission):
    def has_permission(self, request, view):
        return request.user.status == "Student" or "Librarian"


class NotStudentPermission(BasePermission):
    def has_permission(self, request, view):
        return not request.user == "Student"
