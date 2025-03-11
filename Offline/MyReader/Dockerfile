# Use a base Python image or any other image
FROM python:3.12

# Install Git
RUN apt update && \
    apt install -y git cmake make

# Set the working directory

WORKDIR /app

COPY . .
COPY build.sh build.sh

RUN ./build.sh

# Copy the requirements file and install dependencies
RUN pip install --no-cache-dir -r requirements.txt

#ENV MYTEMPO_READER_TYPE="chafon"
ENV LD_LIBRARY_PATH=/usr/local/lib/chafon:/usr/local/lib/assertutils:/usr/local/lib/chafon_rfid:$LD_LIBRARY_PATH

COPY monitor /usr/local/bin/monitor

# Run the Python script
CMD ["/usr/local/bin/monitor"]

