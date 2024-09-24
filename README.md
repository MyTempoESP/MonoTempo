# MonoTempo

Centralized stuff

# Tasks

- [X] Configure network automatically ( not quite there yet,
      solution is currently a workaround with netmon ).

## Next minor version

## IMPORTANT

- [ ] Diagnose/Fix impinj reader
- [ ] Reenvio ( Present first functioning draft )
      deadline: 14d ( fim do mes de set )
- [ ] Exportar banco de dados na máquina local.

### In progress

- [ ] Hotspot
- [ ] Reenvio

### Assigned

- [ ] Install script
- [ ] Start container and stuff automatically at boot
- [ ] Configure impinj reader's time automatically <- FIXME
- [ ] Configure docker-compose.yml to the correct reader
- [ ] Get equipment checkpoint data.

## Per-topic

### Datemon

- [ ] Move all the logic from datemon to MyReader or Envio
      a shell script shouldn't be a container, nor should
      it deal with it's own syncronization.

### Outro

Nice font recommendation: Recursive mono.

### Scratch

```sql
SELECT
    athlete_num,
    MAX(athlete_time),
    checkpoint_id,
    antenna,
    staff,
    tracks.event_id,
    tracks.id
FROM
    athletes_times
JOIN athletes ON athlete_num = athletes.num
JOIN tracks ON tracks.id = athletes.track_id
WHERE
    TIME(athlete_time) > TIME(tracks.largada)
GROUP BY athlete_num, checkpoint_id, antenna, staff, tracks.event_id, tracks.id;
```

