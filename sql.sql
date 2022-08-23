SELECT d.date,d.user, d.station,d.items, e.hours FROM
(SELECT a.date,a.user,a.station,c.usercode,sum(b.items_total) items FROM
	(SELECT ordernum, station, user, DATE(scans.time) as date from scans group by ordernum, station, user, DATE(scans.time)) a INNER JOIN
    (SELECT id, items_total from orders) b on a.ordernum = b.id LEFT JOIN
    (SELECT usercode,username from users) c on a.user = c.username
GROUP BY a.date,a.user,a.station,c.usercode) d LEFT JOIN
(SELECT DATE(clock_in) clockin, payroll_id, sum(paid_hours) hours from shifts group by DATE(clock_in),payroll_id) e on d.date = e.clockin and d.usercode = e.payroll_id
ORDER BY 1,2,3;

//pick/ship list
select a.id,a.date_created,a.statusid,b.ordernum,b.station,b.user,b.time,c.ordernum,c.station,c.user,c.time from orders a
LEFT JOIN (select * FROM scans where station="pick") b ON a.id = b.ordernum
LEFT JOIN (select * FROM scans where station="ship") c ON a.id = c.ordernum
WHERE a.statusid not in (0)
order by 1,5;

//Efficiency
SELECT d.date,d.user, d.items, e.hours, d.items/e.hours FROM
(SELECT a.date,a.user,c.usercode,sum(b.items_total) items FROM
	(SELECT ordernum, station, user, DATE(scans.time) as date from scans where station="pick" group by ordernum, station, user, DATE(scans.time)) a INNER JOIN
    (SELECT id, items_total from orders) b on a.ordernum = b.id LEFT JOIN
    (SELECT usercode,username from users) c on a.user = c.username
GROUP BY a.date,a.user,c.usercode) d LEFT JOIN
(SELECT DATE(clock_in) clockin, payroll_id, sum(paid_hours) hours from shifts group by DATE(clock_in),payroll_id) e on d.date = e.clockin and d.usercode = e.payroll_id
ORDER BY 1,2;

//Group Efficiency
SELECT shipments.date, items/hours efficiency FROM (
select CAST(c.time as date) date, sum(a.items_total) items from orders a
LEFT JOIN (select * FROM scans where station="ship") c ON a.id = c.ordernum
WHERE a.statusid not in (0) and c.time is not null
GROUP BY CAST(c.time as date)
) shipments
LEFT JOIN (select cast(clock_in as date) date,sum(paid_hours) hours FROM shifts WHERE role = 'Shipping' group by cast(clock_in as date)) d on d.date = shipments.date
order by 1;
