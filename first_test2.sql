select titles.title,
	   employees.first_name as 'first name',
       employees.last_name as 'last name',
	   salaries.salary as 'current salary'
from dept_manager, titles, employees, salaries
where dept_manager.to_date >= CURDATE() and
		dept_manager.emp_no = titles.emp_no and
        dept_manager.emp_no = employees.emp_no and
        dept_manager.emp_no = salaries.emp_no and
        salaries.to_date >= curdate() and
        titles.to_date >= curdate();
