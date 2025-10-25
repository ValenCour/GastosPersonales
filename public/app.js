const API_URL_USUARIOS = 'http://localhost:8080/usuarios';
const API_URL_GASTOS = 'http://localhost:8080/gastos';
const userSelect = document.getElementById('user-select');
const gastosList = document.getElementById('gastos-list');
const addUserForm = document.getElementById('add-user-form');
const addGastoForm = document.getElementById('add-gasto-form');

let usuarios = [];
let gastos = [];
let selectedUserId = null;

document.addEventListener('DOMContentLoaded', () => {
    const renderUserOptions = () => {
        userSelect.innerHTML = '<option value="">Seleccione un usuario</option>';
        usuarios.forEach(user => {
            const option = document.createElement('option');
            option.value = user.id_usuario;
            option.textContent = user.nombre_usuario;
            userSelect.appendChild(option);
        });
    };

    const fetchUsuarios = async () => {
        try {
            const response = await fetch(API_URL_USUARIOS);
            if (!response.ok) throw new Error('Error al obtener usuarios');
            usuarios = await response.json();
            renderUserOptions();
        } catch (error) {
            console.error(error);
        }
    };

    userSelect.addEventListener('change', () => {
        selectedUserId = userSelect.value;
        renderGastos();
    });

    addUserForm.addEventListener('submit', async (e) => {
        e.preventDefault(); // Evita que el formulario recargue la página
        
        const newUser = {
            nombre_usuario: document.getElementById('name').value,
            email: document.getElementById('email').value,
            contraseña: document.getElementById('password').value
        };

        try {
            const response = await fetch(API_URL_USUARIOS, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(newUser)
            });

            if (response.status === 201) {
                addUserForm.reset();
                fetchUsuarios(); // Actualiza la lista de usuarios en el dropdown
            } else {
                throw new Error('Error al crear el usuario');
            }
        } catch (error) {
            console.error(error);
        }
    });

    const init = () => {
        fetchUsuarios();
        fetchGastos();
    };

    init();
});