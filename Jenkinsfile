pipeline {
    options {
        buildDiscarder logRotator(artifactNumToKeepStr: '10')
    }
    agent {
        label 'go'
    }
    stages {
        stage('Prepare') {
            steps {
                checkout scm
                sh 'go get github.com/mitchellh/gox'
                sh 'go get -u -v .'
                sh 'mkdir -p builds'
            }
        }
        stage('Build') {
            environment {
                CGO_ENABLED = '0'
            }
            steps {
                sh 'gox -parallel=2 -ldflags="-s -w" -output="builds/{{.Dir}}_{{.OS}}_{{.Arch}}" ./cmd/balloond'
            }
        }
        stage('Cleanup') {
            steps {
                archiveArtifacts 'builds/*'
            }
        }
    }
}
