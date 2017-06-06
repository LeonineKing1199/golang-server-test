(function() {
    'use strict';

    const method   = 'POST';
    const url      = 'http://localhost:8080/users';
    const isAsync  = true;

    const sendRequest = (xhr, firstName, lastName, cb) => {
        xhr.onreadystatechange = cb;
        xhr.open(method, url, isAsync);
        xhr.send({ firstName, lastName });
    };

    const handleBtnSubmit = () => {
        const firstNameDomNode = document.getElementById('first-name-input');
        const lastNameDomNode  = document.getElementById('last-name-input');     

        const firstName = firstNameDomNode.value;
        const lastName  = lastNameDomNode.value;

        const xhr = new XMLHttpRequest();
        sendRequest(xhr, firstName, lastName, () => {
            if (xhr.readyState == XMLHttpRequest.DONE) {
                console.log(`request complated with: ${xhr.status}`);
            }
        });
    };

    document.addEventListener('DOMContentLoaded', () => {
        const addUserBtn = document.getElementById('add-user-btn');

        addUserBtn.addEventListener('mouseup', handleBtnSubmit);
    });
})();