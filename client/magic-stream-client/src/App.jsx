import Home from "./components/home/Home";
import Header from "./components/header/Header";
import { useState, useEffect } from 'react';
import Register from "./components/register/Register";
import Login from "./components/login/Login";
import './App.css'
import {Route, Routes, useNavigate} from 'react-router-dom'

function App() {
  return (
    <>
      <Header />
      <Routes >
        <Route path="/" element={<Home />}></Route>
        <Route path="/register" element={<Register/>}></Route>
        <Route path="/login" element={<Login/>}></Route>
        
      </Routes>
    </>
  );
}

export default App;
