// firebase.ts
import { initializeApp } from "firebase/app";
import { getFirestore } from "firebase/firestore";
import { getAnalytics } from "firebase/analytics";

const firebaseConfig = {
    apiKey: "AIzaSyAbBWyYjc8QavNlYg5TBj0v9lk_7orElKI",
    authDomain: "fuel-data-2e457.firebaseapp.com",
    projectId: "fuel-data-2e457",
    storageBucket: "fuel-data-2e457.appspot.com",
    messagingSenderId: "112036942136",
    appId: "1:112036942136:web:356213587b6922f964b1cd",
    measurementId: "G-MY7CZYJBK2"
};

// Initialize Firebase
const app = initializeApp(firebaseConfig);

// Initialize Firestore
const firestore = getFirestore(app);

// Initialize Analytics
const analytics = getAnalytics(app);

export { firestore, analytics };
