<p align="center">
  <h1>ZENFLOWS INBOX</h1>

Simple **HTTP inbox** for [interfacer-gui](https://github.com/dyne/interfacer-gui/)

</p>

<div align="center">

# Zenflows INBOX

### Federated simple **inbox** for interfacer-gui

</div>

<p align="center">
  <a href="https://dyne.org">
    <img src="https://files.dyne.org/software_by_dyne.png" width="170">
  </a>
</p>

## Building the digital infrastructure for Fab Cities

<br>
<a href="https://www.interfacerproject.eu/">
  <img alt="Interfacer project" src="https://dyne.org/images/projects/Interfacer_logo_color.png" width="350" />
</a>
<br>

### What is **INTERFACER?**

The goal of the INTERFACER project is to build the open-source digital infrastructure for Fab Cities.

Our vision is to promote a green, resilient, and digitally-based mode of production and consumption that enables the greatest possible sovereignty, empowerment and participation of citizens all over the world.
We want to help Fab Cities to produce everything they consume by 2054 on the basis of collaboratively developed and globally shared data in the commons.

To know more [DOWNLOAD THE WHITEPAPER](https://www.interfacerproject.eu/assets/news/whitepaper/IF-WhitePaper_DigitalInfrastructureForFabCities.pdf)

## Zenflows INBOX Features

-   Federated architecture
-   Distributed over tarantool
-   Solid crypto provided by [zenflows-crypto](https://github.com/dyne/zenflows-crypto.git)

# [LIVE DEMO](https://gateway0.interfacer.dyne.org/)

<br>

<div id="toc">

### 🚩 Table of Contents

-   [💾 Install](#-install)
-   [🎮 Quick start](#-quick-start)
-   [🐋 Docker](#-docker)
-   [📋 Testing](#-testing)
-   [😍 Acknowledgements](#-acknowledgements)
-   [🌐 Links](#-links)
-   [👤 Contributing](#-contributing)
-   [💼 License](#-license)

</div>

---

## 💾 Install

For the deployment see the subdirectory `devops`, there is an `ansible` role. It is also available a `Dockerfile` and a `docker-compose.yml`.

**[🔝 back to top](#toc)**

---

## 🎮 Quick start

All request have to be signed with the private key of the `sender` (in `/send`) or `receiver` (in `/read`) agent with [zenflows-crypto](https://github.com/dyne/zenflows-crypto.git) and the signature has to be put in the HTTP request in the header `zenflows-sign`.

### POST `/send`

Send content to a list of receivers.

**Parameters**

|        Name | Required |  Type  | Description                                                                                                                    |
| ----------: | :------: | :----: | ------------------------------------------------------------------------------------------------------------------------------ |
|    `sender` | required |  ULID  | The `sender` is the ID (of the agent in zenflows as string)                                                                    |
| `receivers` | required | ULID[] | The `receivers` is a list of the IDs of the agent (as strings) that should receive the message.                                |
|   `content` | required |  json  | The `content` is saved as JSON inside a postgresql field, when an agent want to see his messages has to make a call to `read`; |

### POST `/read`

Read content for a specific agent.

**Parameters**

|          Name | Required |  Type   | Description                                                                                                       |
| ------------: | :------: | :-----: | ----------------------------------------------------------------------------------------------------------------- |
|    `receiver` | required |  ULID   | The `receiver` is the ID (as string) of the agent we want to read the messages of                                 |
|  `request_id` | required | number  | `request_id` is a random value, in the response the `inbox` service will put the same value in the `receiver_id`. |
| `only_unread` | optional | boolean | There could be a third field `only_unread` that return only the messages for which the `read` flag is `false`;    |

### POST `/set-read`

Mark a specific content as read or unread.

|         Name | Required |  Type   | Description                                                                                                 |
| -----------: | :------: | :-----: | ----------------------------------------------------------------------------------------------------------- |
| `message_id` | required |  ULID   | The `message` is the ID (as string) of the message we want to set status of                                 |
|   `receiver` | required |  ULID   | The `receiver` is the ID (as string) of the agent we want to set status of                                  |
|       `read` | optional | boolean | The `read` flag tells if the `receiver` has read (or not) the message. (Implement the read/unread feature); |

### POST `/count-unread`

Returns the number of messages with the `read` flag set to false.

|       Name | Required | Type | Description                                                                            |
| ---------: | :------: | :--: | -------------------------------------------------------------------------------------- |
| `receiver` | required | ULID | The `receiver` is the ID (as string) of the agent we want to count the unread messages |

**[🔝 back to top](#toc)**

---

## 🐋 Docker

Please refer to [DOCKER PACKAGES](../../packages)

**[🔝 back to top](#toc)**

---

## 📋 Testing

See subdirectory `examples`

**[🔝 back to top](#toc)**

---

## 😍 Acknowledgements

<a href="https://dyne.org">
  <img src="https://files.dyne.org/software_by_dyne.png" width="222">
</a>

Copyleft (ɔ) 2023 by [Dyne.org](https://www.dyne.org) foundation, Amsterdam

Designed, written and maintained by Alberto Lerda
With contributions of Ennio Donato and Puria Nafisi Azizi

### Funding

#### October 2021 – March 2023

This project was funded by the European Regional Development Fund (ERDF) as part of the INTERFACER project.

<img src="docs/assets/funding/logo_erdf2014.png" alt="EU ERDF Programme Logo" width="480">

#### January 2025 – December 2026

The further development of this project was supported as part of the "Fab City – Decentralized Digital Production for Urban Value Creation" project, funded by dtec.bw and supported by the European Union – NextGenerationEU.

<p>
  <img src="docs/assets/funding/dtec.bw_Logo_RGB_komplett.jpg" alt="dtec.bw logo" width="240">
  <img src="docs/assets/funding/nextgeneu_en.jpg" alt="Supported by the European Union – NextGenerationEU" width="360">
</p>

**[🔝 back to top](#toc)**

---

## 🌐 Links

https://www.interfacer.eu/

https://dyne.org/

**[🔝 back to top](#toc)**

---

## 👤 Contributing

1.  🔀 [FORK IT](../../fork)
2.  Create your feature branch `git checkout -b feature/branch`
3.  Commit your changes `git commit -am 'Add some fooBar'`
4.  Push to the branch `git push origin feature/branch`
5.  Create a new Pull Request
6.  🙏 Thank you

**[🔝 back to top](#toc)**

---

## 💼 License

    Zenflows INBOX - Federated simple **inbox** for interfacer-gui
    Copyleft (ɔ) 2023 Dyne.org foundation

    This program is free software: you can redistribute it and/or modify
    it under the terms of the GNU Affero General Public License as
    published by the Free Software Foundation, either version 3 of the
    License, or (at your option) any later version.

    This program is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU Affero General Public License for more details.

    You should have received a copy of the GNU Affero General Public License
    along with this program.  If not, see <http://www.gnu.org/licenses/>.

**[🔝 back to top](#toc)**
