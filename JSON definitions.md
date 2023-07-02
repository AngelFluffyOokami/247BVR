//Document for self reference.

Kill JSON structure:

_id         string
id          uuid
killer      killer
victim      victim
weapon      int
season      int
time        int64
serverInfo  serverInfo

killer:

ownerId     string
occupants   []string
type        int
team        int
position    ???
velocity    ???

victim: 

ownerId     string
occupants   []string
type        int
team        int
position    ???
velocity    ???

serverInfo:

missionId   string
onlineUsers []Online
timeOfDay   int

Death JSON Structure:

id          uuid
victim      victim
season      int
time        int64
serverInfo  serverInfo

victim:

ownerId     string
occupants   []string
type        int
team        int
position    ???
velocity    ???

serverInfo:

missionId   string
onlineUsers []Online
timeOfDay   int