async function fetchAndLog(url) {
    try {
        const response = await fetch(url);
        if (!response.ok) {
            throw new Error('Network response was not ok ' + response.statusText);
        }
        return await response.json();
    } catch (error) {
        console.error('There has been a problem with your fetch operation:', error, url);
    }
}

async function populateEvents(id) {
    const url = `http://localhost:9090/events?userId=${id}`;
    console.log(url);
    const result = await fetchAndLog(url);
    if (!result) return;

    for (let i = 1; i <= 5; i++) {
        const eventElement = document.getElementById(`event${i}`);
        if (eventElement) {
            if (result[i-1]){
                eventElement.innerText = `${result[i - 1].Name} \n ${result[i - 1].Date} ${result[i - 1].Time}`;
            } else {
                eventElement.innerText =  "";
            }
        }
    }
}

async function createEvent(userId, name, date, time, reminders){
    console.log(userId, name, date, time, reminders);
    const url = `http://localhost:9090/events/create?userId=${userId}&name=${name}
    &date=${date}&time=${time}&reminders=${reminders}`;
    const response = await fetch(url, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: { Name: name, Date: date, Time: time, Reminders: reminders },
    });
    if (!response.ok) {
        throw new Error('Network response was not ok');
    }
    return await response.json();
}

async function submitEventForm(event) {
    event.preventDefault();
    // TODO: Implement TG auth
    const userId = 1;

    const name = document.getElementById('event-name').value;
    const date = document.getElementById('event-date').value;
    const time = document.getElementById('event-time').value;

    let reminderBoxes = document.getElementById("event-reminder-box").children;
    let reminders = [];
    for (let i = 0; i < reminderBoxes.length; i++) {
        let box = reminderBoxes[i].children[0]
        if (box.checked) {
            reminders.push(box.value);
        }
    }

    return await createEvent(userId, name, date, time, reminders);
}
