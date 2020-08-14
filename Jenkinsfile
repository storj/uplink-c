pipeline {
    agent {
        docker {
            label 'main'
            image docker.build("storj-ci", "--pull https://github.com/storj/ci.git").id
            args '-u root:root --cap-add SYS_PTRACE -v "/tmp/gomod":/go/pkg/mod'
        }
    }
    options {
          timeout(time: 26, unit: 'MINUTES')
    }
    environment {
        NPM_CONFIG_CACHE = '/tmp/npm/cache'
        COVERDIR = "${ env.BRANCH_NAME != 'master' ? '' : env.WORKSPACE + '/.build/cover' }"
    }
    stages {
        stage('Build') {
            steps {
                checkout scm

                sh 'mkdir -p .build $COVERDIR'

                sh 'service postgresql start'

                sh 'cockroach start-single-node --insecure --store=\'/tmp/crdb\' --listen-addr=localhost:26257 --http-addr=localhost:8080 --cache 512MiB --max-sql-memory 512MiB --background'
            }
        }

        stage('Verification') {
            parallel {
                stage('Lint') {
                    steps {
                        sh 'check-copyright'
                        sh 'check-large-files'
                        sh 'check-imports ./...'
                        sh 'check-peer-constraints'
                        sh 'storj-protobuf --protoc=$HOME/protoc/bin/protoc lint'
                        sh 'storj-protobuf --protoc=$HOME/protoc/bin/protoc check-lock'
                        sh 'check-atomic-align ./...'
                        sh 'check-errs ./...'

                        sh 'go vet ./...'
                        dir('testsuite') {
                            sh  'go vet ./...'
                        }

                        sh 'staticcheck ./...'
                        sh 'golangci-lint --config /go/ci/.golangci.yml -j=2 run'
                        sh 'go-licenses check ./...'
			sh 'make format-c-check'
                    }
                }

                stage('Build') {
                    steps {
                        sh 'make build'
                    }
                }

                stage('Tests') {
                    environment {
                        COVERFLAGS = "${ env.COVERDIR ? '-coverprofile=' + env.COVERDIR + '/tests.coverprofile -coverpkg=./...' : ''}"
                    }
                    steps {
                        sh 'go test -parallel 4 -p 6 -vet=off $COVERFLAGS -timeout 20m -json -race ./... 2>&1 | tee .build/tests.json | xunit -out .build/tests.xml'
                        // TODO enable this later 
                        // sh 'check-clean-directory'
                    }

                    post {
                        always {
                            sh script: 'cat .build/tests.json | tparse -all -top -slow 100', returnStatus: true
                            archiveArtifacts artifacts: '.build/tests.json'
                            junit '.build/tests.xml'
                        }
                    }
                }

                stage('Testsuite') {
                    environment {
                        STORJ_TEST_COCKROACH = 'cockroach://root@localhost:26257/testcockroach?sslmode=disable'
                        STORJ_TEST_POSTGRES = 'postgres://postgres@localhost/teststorj?sslmode=disable'
                        COVERFLAGS = "${ env.COVERDIR ? '-coverprofile=' + env.COVERDIR + '/testsuite.coverprofile -coverpkg=../...' : ''}"
                    }
                    steps {
                        sh 'cockroach sql --insecure --host=localhost:26257 -e \'create database testcockroach;\''
                        sh 'psql -U postgres -c \'create database teststorj;\''
                        sh 'use-ports -from 1024 -to 10000 &'
                        dir('testsuite'){
                            sh 'go test -parallel 4 -p 6 -vet=off $COVERFLAGS -timeout 20m -json -race ./... 2>&1 | tee ../.build/testsuite.json | xunit -out ../.build/testsuite.xml'
                        }
                        // TODO enable this later 
                        // sh 'check-clean-directory'
                    }

                    post {
                        always {
                            sh script: 'cat .build/testsuite.json | tparse -all -top -slow 100', returnStatus: true
                            archiveArtifacts artifacts: '.build/testsuite.json'
                            junit '.build/testsuite.xml'
                        }
                    }
                }
            }
        }

        stage('Coverage') {
            when { not { environment name: 'COVERDIR', value: '' } }
            steps {
                script {
                    def cleaned = []
                    findFiles(glob: '.build/cover/**.coverprofile').each { file ->
                        sh script: "filter-cover-profile < ${file.path} > ${file.path}.clean", returnStatus: true
                        cleaned.push(file.path + '.clean')
                    }
                    sh script: "gocov convert ${cleaned.join(' ')} > .build/cover/combined.json", returnStatus: true
                    sh script: "gocov-xml  < .build/cover/combined.json > .build/cover/combined.xml", returnStatus: true
                    cobertura coberturaReportFile: ".build/cover/combined.xml"
                }
            }
        }
    }

    post {
        always {
            sh "chmod -R 777 ." // ensure Jenkins agent can delete the working directory
            deleteDir()
        }
    }
}
