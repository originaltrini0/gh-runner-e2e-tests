FROM alpine:latest
RUN echo "Hello from container build $(uname -m)" > /test.txt
RUN cat /test.txt
CMD ["cat", "/test.txt"]
