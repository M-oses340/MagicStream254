import Button from 'react-bootstrap/Button'
import Container from 'react-bootstrap/Container'
import Navbar from 'react-bootstrap/Navbar'
import Nav from 'react-bootstrap/Nav'
import {useNavigate, NavLink, Link} from 'react-router-dom'

const Header = () => {
    const navigate = useNavigate();

    return(
        <Navbar bg="dark" variant='dark' expand="lg" sticky="top" className="shadow-sm" >
            <Container>
                <Navbar.Brand>
                    <img
                        alt=""
                        src={logo}
                        width="30"
                        height="30"
                        className="d-inline-block align-top me-2"
                    />
                    Magic Stream
                </Navbar.Brand>
                <Navbar.Toggle aria-controls="main-navbar-nav" />
                <Navbar.Collapse>
                    <Nav className ="me-auto">
                        <Nav.Link as = {NavLink} to="/">
                            Home
                        </Nav.Link>
                        <Nav.Link as = {NavLink} to="/recommended">
                            Recommended
                        </Nav.Link>
                    </Nav>
                </Navbar.Collapse>

            </Container>
        </Navbar>
        
    )
} 