def projectName = 'rightstuff-176212';
def imageName = "gcr.io/${projectName}/jenkins-slave:golang.master";
def feSvcName = "inv-backend"
def baseImageTag = "gcr.io/${projectName}/inv-api:${env.BRANCH_NAME.replace("/", "-")}"
def imageTag = "${baseImageTag}.${env.BUILD_NUMBER}"

podTemplate(cloud: 'local cluster', label: 'golang-k8s', 
    containers: [containerTemplate(name: 'golang', image: imageName, ttyEnabled: true, command: 'cat', alwaysPullImage: true)],
    volumes: [
            hostPathVolume(hostPath: '/var/run/docker.sock', mountPath: '/var/run/docker.sock'),
    ]
) {
    node('golang-k8s') {
        container('golang') {
            checkout scm

            ansiColor('xterm') {
                stage('Build') {
                    sh 'CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o inv-api .'
                }

                stage('Test') {
                    sh 'go test'
                }

                stage('Build docker') {
                    sh "DOCKER_API_VERSION=1.23 docker build -t ${imageTag} ."
                    sh "DOCKER_API_VERSION=1.23 gcloud docker -- push ${imageTag}"
                }

                stage('Deploy') {
                    switch(env.BRANCH_NAME) {
                        // Roll out to canary environment
                        case "canary":
                            sh "sed -i.bak 's#gcr.io/${projectName}/inv-api:1.0.0#${imageTag}#' ./k8s/canary/*.yaml"
                            sh "kubectl --namespace=production apply -f k8s/services/"
                            sh "kubectl --namespace=production apply -f k8s/canary/"
                            sh "echo http://`kubectl --namespace=production get service/${feSvcName} --output=json | jq -r '.status.loadBalancer.ingress[0].ip'` > ${feSvcName}"
                            break
                        // Roll out to production environment
                        case "master":
                            sh "sed -i.bak 's#gcr.io/${projectName}/inv-api:1.0.0#${imageTag}#' ./k8s/production/*.yaml"
                            sh "kubectl --namespace=production apply -f k8s/services/"
                            sh "kubectl --namespace=production apply -f k8s/production/"
                            sh "echo http://`kubectl --namespace=production get service/${feSvcName} --output=json | jq -r '.status.loadBalancer.ingress[0].ip'` > ${feSvcName}"
                            break
                        default:
                            // Create namespace if it doesn't exist
                            sh "kubectl get ns ${env.BRANCH_NAME} || kubectl create ns ${env.BRANCH_NAME}"
                            // Don't use public load balancing for development branches
                            sh "sed -i.bak 's#LoadBalancer#ClusterIP#' ./k8s/services/frontend.yaml"
                            sh "sed -i.bak 's#gcr.io/${projectName}/inv-api:1.0.0#${imageTag}#' ./k8s/develop/*.yaml"
                            sh("kubectl --namespace=${env.BRANCH_NAME} apply -f k8s/services/")
                            sh("kubectl --namespace=${env.BRANCH_NAME} apply -f k8s/develop/")
                            echo 'To access your environment run `kubectl proxy`'
                            echo "Then access your service via http://localhost:8001/api/v1/proxy/namespaces/${env.BRANCH_NAME}/services/${feSvcName}:80/"
                            break
                    }
                }
            }
        }
    }
}