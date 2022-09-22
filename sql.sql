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

//error reporting
select a.orderid, a.issue,b.user,b.time from errors a inner join scans b on a.orderid = b.ordernum where b.station="pick" and a.issue in ('Incorrect','Missing')order by 1 desc;

//errors per hour
select user, round(errors/hours,3) error_rate FROM (select user,usercode,count(*) as errors FROM (select a.orderid, a.issue,b.user,c.usercode,b.time from errors a inner join scans b on a.orderid = b.ordernum left join users c on b.user=c.username where b.station="pick" and a.issue in ('Incorrect','Missing') and time between '2022-08-15' and '2022-08-20') d GROUP BY user,usercode) e LEFT JOIN (select payroll_id, sum(scheduled_hours) hours FROM shifts where clock_in between '2022-08-15' and '2022-08-20' group by payroll_id) f on e.usercode = f.payroll_id

//Recent errors
select a.orderid, b.user, a.issue, a.comment, b.time FROM errors a left join scans b on a.orderid = b.ordernum WHERE a.issue in ('Missing','Incorrect') AND b.station = 'pick' order by 5 desc limit 10

//Hours worked
select b.username,a.ratio FROM (
SELECT payroll_id, sum(paid_hours)/sum(scheduled_hours) as ratio
FROM shifts
where cast(clock_in as date) between '2022-08-01' and '2022-08-31'
group by payroll_id
) a
INNER JOIN users b on a.payroll_id = b.usercode
//Service Level
SELECT week, sum(case when SL < 3 then 1 else 0 end)/count(*) as SL, sum(case when SL < 4 then 1 else 0 end)/count(*) as SL1 FROM (select DATE_ADD(cast(a.date_created as date), INTERVAL(-WEEKDAY(cast(a.date_created as date))) DAY) as week,TOTAL_WEEKDAYS(b.time,a.date_created) - 1 as SL FROM orders a LEFT JOIN scans b ON a.id = b.ordernum where b.station = 'ship') c GROUP BY week ORDER BY 1
