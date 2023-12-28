from django.db import models
from books.models import Book
from users.models import User


class Review(models.Model):
    GRADES = (
        (1, 1), (2, 2), (3, 3), (4, 4), (5, 5)
    )
    author = models.ForeignKey(User, on_delete=models.CASCADE)
    book = models.ForeignKey(Book, on_delete=models.CASCADE)
    text = models.TextField(default="", blank=True)
    grade = models.IntegerField(choices=GRADES, null=True)
    created_time = models.DateTimeField(auto_now_add=True)

    class Meta:
        ordering = ["-created_time"]
        db_table = "reviews"
        verbose_name = "Отзыв"
        verbose_name_plural = "Отзывы"

    def save(self, *args, **kwargs):
        if self.id:
            return

        super(Review, self).save(*args, **kwargs)

        self.book.total_rating += self.grade
        self.book.reviews_quantity += 1

        if self.book.reviews_quantity > 0:
            self.book.rating = round(float(self.book.total_rating)
                                / float(self.book.reviews_quantity), 2)
        else:
            self.book.rating = 0

        self.book.save()

    def __str__(self):
        return f"Отзыв от {self.author} на книгу {self.book}"
