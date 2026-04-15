pipeline {
    agent any                        // Di mana pipeline dijalankan

    environment {                    // Variabel lingkungan
        APP_ENV = 'staging'
    }

    tools {                          // Tools yang digunakan (Maven, JDK, dsb.)
        maven 'Maven 3.8'
    }

    stages {
        stage('Checkout') { ... }    // Ambil kode dari repository
        stage('Build')    { ... }    // Kompilasi / build
        stage('Test')     { ... }    // Jalankan pengujian
        stage('Analyze')  { ... }    // Analisis SonarQube
        stage('Deploy')   { ... }    // Deployment
    }

    post {                           // Aksi setelah pipeline selesai
        always   { echo 'Pipeline selesai' }
        success  { echo 'Pipeline berhasil!' }
        failure  { echo 'Pipeline gagal!' }
    }
}