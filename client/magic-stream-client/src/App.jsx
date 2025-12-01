import Home from "./components/home/Home";
import Header from "./components/header/Header"
import { useState, useEffect } from 'react'
import './App.css'

function App() {
  return (
    <>
      <Header/>
      <Home />
    </>
  );
}

export default App;
