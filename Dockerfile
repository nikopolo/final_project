FROM ubuntu:latest

WORKDIR /app

COPY task_scheduler ./

COPY web/ web/

COPY scheduler.db ./

EXPOSE 7540

CMD ["./task_scheduler"]