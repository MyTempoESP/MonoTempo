# MonoTempo

container configuration for monitoring:

|      Container        |                                        Description                                         |
| --------------------- | ------------------------------------------------------------------------------------------ |
| AA2                   | An Arduino controlling an LCD screen ( Go/C/forth )                                        |
| mytempo-chafon-reader | My python library for controlling an RFID module ( Python/C/C++ )                          |
| MyReader              | A wrapper around the reader for sending tags through a Message Broker ( RabbitMQ/Python )  |
| Envio                 | Simple program to manage multiple sqlite databases and store tags ( RabbitMQ/Go/Sqlite )   |
| Receba                | Fetch equipment data from the mytempo API ( PHP/Go/Sqlite )                                |
| Reenvio               | Process/Send stored tags to the mytempo API ( PHP/Go/Sqlite )                              |

