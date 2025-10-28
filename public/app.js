const API_URL_USUARIOS = 'http://localhost:8080/usuarios';
const API_URL_GASTOS = 'http://localhost:8080/gastos';
const userSelect = document.getElementById('user-select');
const gastosList = document.getElementById('gastos-list');
const addUserForm = document.getElementById('add-user-form');
const addGastoForm = document.getElementById('add-gasto-form');

let usuarios = [];
let gastos = [];
let id_seleccionado = -1;

document.addEventListener('DOMContentLoaded', () => {
    const renderUserOptions = () => {
        userSelect.innerHTML = '<option value="-1">Seleccione un usuario</option>';
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

    const fetchGastosPorId = async (id) => {
        gastos = [];
        try {
            const response = await fetch(API_URL_USUARIOS + "/" + id + "/" + "gastos");
            if (!response.ok) throw new Error('Error al obtener gastos');
            gastos = await response.json();
        } catch (error) {
            console.error(error);
        }
        renderGastosId(id);
    };

    const renderGastosId = (id) => {
        gastosList.innerHTML = ''; 
        if (id === -1) {
            gastosList.innerHTML = '<p>Seleccione un usuario para ver sus gastos.</p>';
            return;
        }

        if (gastos.length === 0) {
            gastosList.innerHTML = '<p>Este usuario no tiene gastos registrados.</p>';
            return;
        }

        gastosList.innerHTML = `
        <div class="gastos-header">
            <span>Categoria</span>
            <span>Monto</span>
            <span>Fecha</span>
            <span></span>
        </div>
        `;

        gastos.forEach(gasto => {
            const li = document.createElement('li');
            li.innerHTML = `
                <span>${gasto.categoria}</span>
                <span>$${gasto.monto}</span>
                <span>${new Date(gasto.fecha).toLocaleString()}</span>
                <button class="delete-btn" data-id="${gasto.id_gasto}">Eliminar</button>
            `;
            gastosList.appendChild(li);
        });
    };

    userSelect.addEventListener('change', () => {
        id_seleccionado = userSelect.value;
        fetchGastosPorId(id_seleccionado);
    });

    addUserForm.addEventListener('submit', async (e) => {
        e.preventDefault();
        
        const nuevo_usuario = {
            nombre_usuario: document.getElementById('name').value,
            email: document.getElementById('email').value,
            contraseña: document.getElementById('password').value
        };

        try {
            const response = await fetch(API_URL_USUARIOS, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(nuevo_usuario)
            });

            if (response.status === 201) {
                addUserForm.reset();
                fetchUsuarios();
            } else {
                throw new Error('Error al crear el usuario');
            }
        } catch (error) {
            console.error(error);
        }
    });

    addGastoForm.addEventListener('submit', async (e) => {
        e.preventDefault();
        if(id_seleccionado !== -1){
            const nuevo_gasto = {
                id_usuario:     parseInt(id_seleccionado),
                monto:          document.getElementById('monto').value,
                medio_de_pago:  document.getElementById('medio_de_pago').value,
                fecha:          new Date(document.getElementById('fecha').value).toISOString(),
                categoria:      document.getElementById('categoria').value
            }

            try {
                const response = await fetch(API_URL_GASTOS, {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify(nuevo_gasto)
                });

                if (response.status === 201) {
                    addGastoForm.reset();
                    fetchGastosPorId(id_seleccionado);
                } else {
                    throw new Error('Error al crear el usuario');
                }
            } catch (error) {
                console.error(error);
            }
        }
    })

    gastosList.addEventListener('click', async (e) => {
        if (e.target.classList.contains('delete-btn')) {
            gasto_id = e.target.getAttribute('data-id');
            try {
                const response = await fetch(API_URL_GASTOS + "/" + gasto_id, {method: 'DELETE'});
                if (response.status != 204) throw new Error('Error al obtener gastos');
                gastos = await response.json();
            } catch (error) {
                console.error(error);
            }
            fetchGastosPorId(id_seleccionado);
        }
    })

    const init = () => {
        fetchUsuarios();
        fetchGastosPorId(id_seleccionado);
        addUserForm.reset();
        addGastoForm.reset();
    };

    init();
});