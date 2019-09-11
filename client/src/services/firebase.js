import firebase from "firebase";

var firebaseConfig = {
  apiKey: "AIzaSyDsrSK3MKZabXFhYXES7MD5NCpyKvhejfw",
  authDomain: "controle-acesso-port.firebaseapp.com",
  databaseURL: "https://controle-acesso-port.firebaseio.com",
  projectId: "controle-acesso-port",
  storageBucket: "controle-acesso-port.appspot.com",
  messagingSenderId: "996684619286",
  appId: "1:996684619286:web:d65e2d70d70f52cb"
};
// Initialize Firebase
var app = firebase.initializeApp(firebaseConfig);
var db = firebase.firestore(app);

const teste = db
  .collection("registro-veiculo")
  .get()
  .then(function(querySnapshot) {
    querySnapshot.forEach(function(doc) {
      // doc.data() is never undefined for query doc snapshots
      console.log(doc.id, " => ", doc.data());
    });
  });
console.log("teste1", teste);

var dados = teste.data();

console.log("teste2", dados);
export default teste;
