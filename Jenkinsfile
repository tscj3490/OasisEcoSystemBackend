node {
  def IMAGE_NAME
  def DOCKER_REGISTRY
  def OC_HOST
  def K8S_CONFIG_ENC=""
  def K8S_NAMESPACE=""

  switch (env.BRANCH_NAME) {

    case "master":
      K8S_NAMESPACE="oasis"
      K8S_CONFIG_ENC="k8s.merkins.io.conf.enc"
      K8S_TEMPLATE_VALUES="values.pro.yml"
    break;
    case "develop":
      K8S_NAMESPACE="oasis"
      K8S_CONFIG_ENC="k8s2.merkins.io.conf.enc"
      K8S_TEMPLATE_VALUES="values.yml"
    break;
    default:
      echo "Branch ignored $env.BRANCH_NAME"
  }  

  if (env.K8S_NAMESPACE!=""){
    stage ("Prepare environment") {
      checkout scm

      docker.image('jhidalgo3/kubernetes-node-docker:7-alpine').pull()
      docker.image('jhidalgo3/kubernetes-node-docker:7-alpine').inside('-u root') {
          stage ("Calculate variables") {
              IMAGE_NAME = sh(returnStdout: true, script: "echo {{.docker_registry}}/{{.app_name}}/{{.BRANCH_NAME}}:{{commit}} | dante-cli --values ${K8S_TEMPLATE_VALUES}").trim()
              DOCKER_REGISTRY =sh(returnStdout: true, script: "echo {{.docker_registry}} | dante-cli --values ${K8S_TEMPLATE_VALUES}").trim()
              OC_HOST = sh(returnStdout: true, script: "echo {{.oc_host}} | dante-cli --values ${K8S_TEMPLATE_VALUES}").trim()
              sh "mkdir distYaml | true"
              sh "rm distYaml/* | true"
              sh "dante-cli --values ${K8S_TEMPLATE_VALUES} --template k8s:distYaml"
          }
        }
        
        stage ("Create docker Image"){
          withCredentials([[$class: "UsernamePasswordMultiBinding", usernameVariable: 'DOCKERHUB_USER', passwordVariable: 'DOCKERHUB_PASS', credentialsId: 'registry.jhidalgo3.me']]) {
            sh "docker login --username $DOCKERHUB_USER --password '$DOCKERHUB_PASS' $DOCKER_REGISTRY"
          }

          sh "docker build --rm -t ${IMAGE_NAME} ."
          sh "docker push ${IMAGE_NAME}"
        }

      docker.image('jhidalgo3/kubernetes-node-docker:7-alpine').inside('-u root') {    
        stage ("Deploy Kubernetes"){
          sh "mkdir ~/.kube"

          withCredentials([string(credentialsId: 'MASTER_KEY', variable: 'MASTER_KEY')]){
            sh "openssl enc -aes-256-cbc -d -in ${K8S_CONFIG_ENC} -out ~/.kube/config -k ${MASTER_KEY}"
          }
          sh "kubectl apply -f distYaml -n ${K8S_NAMESPACE}"
        }
      }
    }
  }  
}
