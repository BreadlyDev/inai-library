from django.db import models
from users.models import User


class Message(models.Model):
    author = models.ForeignKey(User, on_delete=models.CASCADE, related_name='authored_messages')
    title = models.CharField(max_length=150)
    text = models.TextField()
    recipient = models.ForeignKey(User, on_delete=models.CASCADE, related_name='received_messages')
    created_time = models.DateTimeField(auto_now_add=True)

    class Meta:
        ordering = ["-created_time"]
        db_table = "messages"

    def __str__(self):
        return f"Message by {self.author} to {self.recipient}"
