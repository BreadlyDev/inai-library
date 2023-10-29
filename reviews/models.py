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

    def save(self, *args, **kwargs):
        super(Review, self).save(*args, **kwargs)

        book = self.book
        reviews = Review.objects.filter(books_id=book)
        total_rating = sum([review.grade for review in reviews if review.grade is not None])
        num_reviews = len(reviews)

        if num_reviews > 0:
            book.rating = round(total_rating / num_reviews, 2)
        else:
            book.rating = 0

        book.save()

    def __str__(self):
        return f"Review by {self.author} on {self.book} book"
