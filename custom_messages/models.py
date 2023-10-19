from django.db import models
from accounts.models import User


class Message(models.Model):
    author = models.ForeignKey(User, on_delete=models.CASCADE, related_name='authored_messages')
    title = models.CharField(max_length=150)
    text = models.TextField()
    recipient = models.ForeignKey(User, on_delete=models.CASCADE, related_name='received_messages')
    created_time = models.DateTimeField(auto_now_add=True)

    def __str__(self):
        return f"Message from {self.author} to {self.recipient}"