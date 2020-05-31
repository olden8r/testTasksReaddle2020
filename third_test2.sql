select departments.dept_name as 'department name'
		, count(distinct dept_emp.emp_no) as 'current employee count'
        , sum(salaries.salary) as 'current sum salary'
from dept_emp, departments, salaries
where departments.dept_no = dept_emp.dept_no
		and dept_emp.emp_no = salaries.emp_no
		and dept_emp.to_date >= curdate() 
		group by departments.dept_no ;