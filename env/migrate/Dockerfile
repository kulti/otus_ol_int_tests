FROM migrate/migrate

COPY migrate.sh /app/migrate.sh
COPY migrations /app/migrations

ENTRYPOINT [ "/app/migrate.sh" ]
CMD [ "up" ]
