FROM scratch

COPY kvs .

# Create a directory in the container to store the certificates
ADD certificates /certificates

# Copy the entire certificates folder from the host into the container
# COPY certificates /certificates

EXPOSE 8080

CMD ["/kvs"]