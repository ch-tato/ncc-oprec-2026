# Laporan Penugasan Modul 1 Open Recruitment NCC 2026

|Nama|NRP|
|---|---|
|Muhammad Quthbi Danish Abqori|5025241036|

- **URL API (Endpoint `/health`):** `http://103.181.182.109:8080/health`
- **Github Repository:** `https://github.com/ch-tato/ncc-oprec-2026/tree/module-1`

#### 1. Deskripsi Singkat Service
Service ini adalah sebuah *microservice* HTTP ringan yang dikembangkan menggunakan bahasa **Golang** (menggunakan *standard library* `net/http` tanpa *framework* eksternal). Tujuan utama pemilihan Golang adalah untuk menghasilkan *statically linked binary* yang sangat optimal dan aman ketika di-*dockerization*.

Berikut adalah snippet kode dari [`main.go`](./main.go):

![Snippet of main.go](<img/Screenshot from 2026-04-14 23-49-07.png>)

#### 2. Penjelasan Endpoint `/health`
Endpoint `/health` diatur untuk merespons HTTP *request* dengan *method* GET. Endpoint ini mengembalikan HTTP Status 200 OK beserta *response body* dalam format JSON. 
Selain memberikan status *service*, endpoint ini juga dikustomisasi untuk memberikan informasi *real-time* berupa:
* Nama & NRP
* Waktu saat ini (*Timestamp* format RFC3339)
* *Uptime* (lama waktu *service* telah berjalan sejak pertama kali dihidupkan)

![JSON Response](<img/Screenshot from 2026-04-14 23-47-51.png>)

#### 3. Bukti Endpoint Dapat Diakses

![curl test](<img/Screenshot from 2026-04-14 23-46-36.png>)
![bruno test](<img/Screenshot from 2026-04-14 23-36-08.png>)

#### 4. Penjelasan Proses Build dan Run Docker
Proses *build* dirancang dengan fokus pada efisiensi ukuran, keamanan, dan *best practice*:
* **Multi-stage Build & Optimasi Ukuran:** Menggunakan dua *stage*. *Stage* pertama (`golang:1.26-alpine`) bertugas melakukan kompilasi kode menjadi *static binary*. *Stage* kedua (`alpine:3.23.3`) hanya mengambil hasil *binary* tersebut. Ini mengoptimasi ukuran *image* secara drastis (dari ratusan MB menjadi <15MB).
* **Environment Variable:** Konfigurasi *port* tidak di-*hardcode*, melainkan di-*inject* melalui instruksi `ENV` di [`Dockerfile`](./Dockerfile) dan [`docker-compose.yml`](./docker-compose.yml).
* **.dockerignore:** Digunakan untuk memfilter direktori `.git` dan `.env` agar tidak masuk ke dalam *build context*, sehingga proses *build* lebih cepat dan *credential* tetap aman.
* **Instruction HEALTHCHECK:** Ditambahkan langsung ke dalam *Dockerfile* (`HEALTHCHECK --interval=30s ... CMD wget ...`) untuk memonitor status internal aplikasi secara berkala dari dalam *container*.
* **Security:** Menambahkan konfigurasi `adduser` untuk menjalankan aplikasi menggunakan *user* non-root (`appuser`).

![building docker image](<img/Screenshot from 2026-04-14 21-48-10.png>)
![built container on the docker desktop](<img/Screenshot from 2026-04-14 21-47-21.png>)

#### 5. Penjelasan Proses Deployment ke VPS
Deployment ke Virtual Machine dilakukan melalui protokol SSH. Proses *deployment* dan orkestrasi *container* menggunakan **Docker Compose** dengan detail berikut:
* **Install Docker & Docker Compose:** Pastikan Docker dan Docker Compose sudah terinstall di VPS.

![installing docker](<img/Screenshot from 2026-04-14 23-28-16.png>)

* **Transfer Kode:** Menggunakan Git untuk melakukan *cloning repository* secara aman ke dalam lingkungan VPS.

![cloning repo](<img/Screenshot from 2026-04-14 23-31-13.png>)

* **Manajemen Rahasia:** Membuat ulang file `.env` secara manual di VPS agar kredensial tidak terekspos di Git.
* **Port Configuration:** Menggunakan pemetaan *port* yang jelas di `docker-compose.yml` (`8080:8080`) untuk menghubungkan *port host* publik dengan *port container* internal.
* **Restart Policy:** Menambahkan instruksi `restart: unless-stopped` agar *service* memiliki resiliensi tinggi dan otomatis hidup kembali jika VPS mengalami *reboot* atau gangguan.

![alt text](<img/Screenshot from 2026-04-14 23-34-39.png>)

#### 6. Kendala yang Dihadapi & Solusi
* **Kendala:** Saat melakukan proses *build* Docker pertama kali di mesin lokal (Linux), terjadi kendala koneksi ke Docker daemon dengan pesan *error*: `failed to connect to the docker API at unix:///home/.../docker.sock`.
* **Solusi:** Melakukan *troubleshooting* dengan mengecek *context* Docker yang aktif. Mengubah *context* kembali ke *default* sistem (`docker context use default`) dan memastikan *service daemon* Docker di latar belakang (*systemctl*) berstatus aktif, sehingga proses *build* dapat dilanjutkan dengan sukses.

![alt text](<img/Screenshot from 2026-04-14 23-59-33.png>)