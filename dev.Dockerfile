# Add wkhtmlpdf image source based on alpine 3.20 with "wkhtmltopdf" as reference name
FROM surnet/alpine-wkhtmltopdf:3.20.0-0.12.6-full as wkhtmltopdf

# Use golang 1.23 on alpine 3.20 as image base
FROM golang:1.23-alpine3.20

# Add wkhtmltopdf required package
RUN apk add --no-cache \
    libstdc++ libx11 libxrender \
    libxext fontconfig freetype \
    ttf-droid ttf-freefont ttf-liberation 

# Copy the wkhtmltopdf binary from "wkhtmltopdf" reference image
COPY --from=wkhtmltopdf /bin/wkhtmltopdf /usr/local/bin/wkhtmltopdf
COPY --from=wkhtmltopdf /bin/wkhtmltoimage /usr/local/bin/wkhtmltoimage
COPY --from=wkhtmltopdf /bin/libwkhtmltox* /usr/local/bin/

# Install golang air autoreload package
RUN go install github.com/air-verse/air@v1.52.3

# Set the workdir to app_src folder, all files in this project directory will be mounted here
WORKDIR /app_src

# Run the autoreload 
CMD ["air", "-c", ".air.toml"]
