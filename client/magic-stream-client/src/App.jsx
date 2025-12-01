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
      <Header handleLogout = {handleLogout}/>
      <Routes path="/" element = {<Layout/>}>
        <Route path="/" element={<Home updateMovieReview={updateMovieReview}/>}></Route>
        <Route path="/register" element={<Register/>}></Route>
        <Route path="/login" element={<Login/>}></Route>
        <Route element = {<RequiredAuth/>}>
            <Route path="/recommended" element={<Recommended/>}></Route>
            <Route path="/review/:imdb_id" element={<Review/>}></Route>
            <Route path="/stream/:yt_id" element={<StreamMovie/>}></Route>
        </Route>
      </Routes>
    </>
  );
}

export default App;
