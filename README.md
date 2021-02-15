# Guest API

<h3>OS independent start guide: </h3>

<ul>
    <li>run docker</li>
    <li>cd api/database </li>
    <li>docker-compose up</li>
    <li>docker exec -it dbTestContainer bash</li>
    <li>mysql --user=root --password="Password123" dbTestContainer < /migrations/migrations/init.sql</li>
    <li>go run api/cmd/main.go</li>
</ul>

<h3>Tests: </h3>
(sorry this is a bit of effort, ran out of time)
<ul>
    <li>run the mock commands found in the make file manually</li>
</ul>