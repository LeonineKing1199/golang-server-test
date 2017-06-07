(function() {
    'use strict';

    const method   = 'POST';
    const url      = 'http://localhost:8080/users';
    const isAsync  = true;

    const makeTableColumn = (textContent) => {
        const td = document.createElement('td');
        td.textContent = textContent;
        return td;
    };

    const makeUserRow = ({ FirstName: firstName, LastName: lastName }) => {
        const tr = document.createElement('tr');

        const fnameTD = makeTableColumn(firstName);
        const lnameTD = makeTableColumn(lastName);

        tr.appendChild(fnameTD);
        tr.appendChild(lnameTD);

        return tr;
    };

    const sendRequest = (xhr, FirstName, LastName, cb) => {
        xhr.onreadystatechange = cb;
        xhr.open(method, url, isAsync);
        xhr.setRequestHeader('Content-type','application/json');
        xhr.send(JSON.stringify({ FirstName, LastName }));
    };

    const handleBtnSubmit = (btn) => () => {
        const firstNameDomNode = document.getElementById('first-name-input');
        const lastNameDomNode  = document.getElementById('last-name-input');     

        const tableBodyDomNode = document.getElementById('user-table-body');

        const firstName = firstNameDomNode.value || '';
        const lastName  = lastNameDomNode.value  || '';

        if (!firstName || !lastName) {
            return;
        }

        const xhr = new XMLHttpRequest();
        sendRequest(xhr, firstName, lastName, () => {
            btn.removeAttribute('disabled');

            if (xhr.readyState == XMLHttpRequest.DONE && xhr.status === 200) {
                firstNameDomNode.value = '';
                lastNameDomNode.value  = '';

                const users = JSON.parse(xhr.responseText);
                
                tableBodyDomNode.innerHTML = '';
                users.forEach((user) => {
                    tableBodyDomNode.appendChild(makeUserRow(user));
                });
            }
        });

        btn.setAttribute('disabled', '');
    };

    document.addEventListener('DOMContentLoaded', () => {
        const addUserBtn = document.getElementById('add-user-btn');

        addUserBtn.addEventListener('mouseup', handleBtnSubmit(addUserBtn));
    });
})();