pipeline {
    agent any
    
    parameters {
        choice(name: 'ARCH', choices: ['arm64', 'amd64'], description: 'Select the target architecture')
        choice(name: 'OS', choices: ['linux', 'darwin', 'windows'], description: 'Select the target operating system')
    }

    environment {
        REPO = 'https://github.com/Andrey-Ilin/kbot'
        BRANCH = 'main'
        TARGETOS = "${params.OS}"
        TARGETARCH = "${params.ARCH}"
    }

    stages {
        
        stage("clone") {
            steps {
                echo 'CLONE THE REPOSITORY'
                git branch: "${BRANCH}", url: "${REPO}"
            }
        }
        stage("test") {
            steps {
               echo 'TEST EXECUTION STARTS'
               sh 'make test'
            }
        }
        stage("build") {
            steps {
               echo 'BUILD STARTS'
               sh 'make build TARGETOS="${TARGETOS}" TARGETARCH="${TARGETARCH}"'
            }
        }
        stage("image") {
            steps {
               echo 'IMAGE CREATION STARTS'
               sh 'make image TARGETOS="${TARGETOS}" TARGETARCH="${TARGETARCH}"'
            }
        }
        stage("push") {
            steps {
                script {
                    docker.withRegistry('', 'dockerhub') {
                        echo 'PUSH STARTS'
                        sh 'make push TARGETOS="${TARGETOS}" TARGETARCH="${TARGETARCH}"'
                    }
                } 
            }
        }
    }
}
