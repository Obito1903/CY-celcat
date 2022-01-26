# Login

Grab the request verification token from here :
https://services-web.u-cergy.fr/calendar/LdapLogin
exemple :
```html
<input name="__RequestVerificationToken" type="hidden" value="CfDJ8O3fKIvt5EFPiHSgETpVhO2aRXYqhYQBRdU9Wfi15C_ahxX3quryziHhBOcDwYTFEOHq7NJ75RVtNFdavu-xeSWUgsCCPtKphhwsX_kpNPUKJNqZQjS5UJon0afxpYigLP-xuMAyz2xy6q6u8Niiq0M" />
```
Save the cookie send to you.

Then send the `Name` , `Password` and `__RequestVerificationToken` to https://services-web.u-cergy.fr/calendar/LdapLogin/Logon with the previous cookie
Save the new cookie sent to you and extract your federation id from the respon header `location`

# Get Data

## List of event

send a form with this data :

|       field       |    data    | desc                                        |
| :---------------: | :--------: | ------------------------------------------- |
|      `start`      | aaaa-mm-dd | Start date of the events to query           |
|       `end`       | aaaa-mm-dd | End date of the events to query             |
|     `resType`     |    xxx     | Type de ressource a query (104 for our use) |
|     `calView`     |   string   | usualy `agendaWeek` or `month`              |
| `federationIds[]` |  22xxxxxx  | Id of the Student                           |

This will send you a json response containing a list of all the event for the time period requested :

```json
{
	"id": "-1347128091:-662573064:1:153781:22",
	"start": "2022-01-24T09:45:00",
	"end": "2022-01-24T13:00:00",
	"allDay": false,
	"description": "TD\r\n\r\n<br />\r\n\r\nLV1 / TOEIC\r\n\r\n<br />\r\n\r\nPAU E101 LABO DE LANGUES 49p<br />PAU E102 SALLE POLYVALENTE (TD ET INFO) 40p\r\n\r\n<br />\r\n\r\nCOYNAULT MAGDALEN<br />GALAN SEBASTIEN\r\n",
	"backgroundColor": "#badefc",
	"textColor": "#000000",
	"department": "D : CY TECH",
	"faculty": null,
	"eventCategory": "TD",
	"sites": [
		"EISTI",
		"EISTI"
	],
	"modules": [
		"DIALVU2D"
	],
	"registerStatus": 2,
	"studentMark": 0,
	"custom1": null,
	"custom2": null,
	"custom3": null
}
```

## Event details

To get mor details about an event send a request here https://services-web.u-cergy.fr/calendar/Home/GetSideBarEvent with the following fields

|   field   |              data              | desc                     |
| :-------: | :----------------------------: | ------------------------ |
| `eventId` | xxxxxxxxxx:xxxxxxxx:x:xxxxx:xx | id of the event to query |

this will return you a json with the following fields :

|     field     |      data       | desc                                        |
| :-----------: | :-------------: | ------------------------------------------- |
| `fereationId` |    22xxxxxx     | The federationId of the event (usualy null) |
| `entityType`  |     integer     |                                             |
|  `elements`   | list of objects | List of elements describing the event       |

### Elements fields

|        field         |  data   | desc                              |
| :------------------: | :-----: | --------------------------------- |
|       `label`        | string  | Name or categorie of the elements |
|      `content`       | string  | content of the element            |
|     `entityType`     | integer | internal category of the element  |
| `isStudentSpecific`  | boolean |                                   |
|    `fereationId`     | integer |                                   |
| `assignmentContext`  | string  |                                   |
| `containsHyperlinks` | boolean |                                   |
|      `isNotes`       | boolean |                                   |


#### Labels list :
| name                          | entityType | decs                                                                                         |
| ----------------------------- | ---------- | -------------------------------------------------------------------------------------------- |
| `Time`                        | 0          |                                                                                              |
| `Catégorie`                   | 0          | Category of the event (CM, TD...)                                                            |
| `Matière`                     | 100        |                                                                                              |
| `Salle` or `Salles`           | 102        | In case of Salles you need to add the next entry with label `null` and entityType `102`      |
| `Enseignant` or `Enseignants` | 101        | In case of Enseignants you need to add the next entry with label `null` and entityType `102` |
| `Notes`                       | 0          |                                                                                              |
| `Name`                        | 0          | Often used to add small info to events like if it is canceled                                |

__Notes__:
Some of those labels might not exist or might have some fields set to `null`
