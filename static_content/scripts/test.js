(function() {
    'use strict';

    const method   = 'POST';
    const url      = 'http://localhost:8080/users';
    const isAsync  = true;

    const sendRequest = (xhr, FirstName, LastName, cb) => {
        xhr.onreadystatechange = cb;
        xhr.open(method, url, isAsync);
        xhr.setRequestHeader('Content-type','application/json');
        xhr.send(JSON.stringify({ FirstName, LastName }));
    };

    const handleBtnSubmit = (btn) => () => {
        const firstNameDomNode = document.getElementById('first-name-input');
        const lastNameDomNode  = document.getElementById('last-name-input');     

        const firstName = firstNameDomNode.value || '';
        const lastName  = lastNameDomNode.value  || '';

        if (!firstName || !lastName) {
            return;
        }

        const xhr = new XMLHttpRequest();
        sendRequest(xhr, firstName, lastName, () => {
            btn.removeAttribute('disabled');

            if (xhr.readyState == XMLHttpRequest.DONE) {
                console.log(`request complated with: ${xhr.status}`);
            }
        });

        btn.setAttribute('disabled', '');
    };

    document.addEventListener('DOMContentLoaded', () => {
        const addUserBtn = document.getElementById('add-user-btn');

        addUserBtn.addEventListener('mouseup', handleBtnSubmit(addUserBtn));
    });
})();